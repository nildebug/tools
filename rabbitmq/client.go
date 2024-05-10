package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type Client struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	connectURL   string
	mode         int
	exchangeName string
	queueName    string
	dataHandler  func(data []byte)
}

func (c *Client) debugString() string {
	return fmt.Sprintf("exchangeName:%s queueName:%s", c.exchangeName, c.queueName)
}

const (
	AMQP_MODE_EXCHANGE = iota + 100 //发布/订阅（Publish/Subscribe）
	AMQP_MODE_QUEUE                 //队列模式（Queue）
	AMQP_MODE_TOPICS                //主题（Topics）
)

type URLConfig struct {
	UserName string
	PassWord string
	IPAddr   string
	Port     int
	Vhost    string
}

func (u *URLConfig) GetDialURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", u.UserName, u.PassWord, u.IPAddr, u.Port, u.Vhost)
}
func connect(ctx context.Context, url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

var isConsumer bool

func NewConsumer(ctx context.Context, config URLConfig, mode int, exchangeName string, queueName string, dataHandler func(data []byte)) (*Client, error) {
	conn, err := connect(ctx, config.GetDialURL())
	if err != nil {
		return nil, err
	}
	if dataHandler == nil {
		return nil, fmt.Errorf("dataHandler is nil")
	}
	client := &Client{
		conn:         conn,
		connectURL:   config.GetDialURL(),
		mode:         mode,
		queueName:    queueName,
		exchangeName: exchangeName,
		dataHandler:  dataHandler,
	}
	isConsumer = true
	if err = client.initConsume(); err != nil {
		return nil, err
	}
	fmt.Println("rabbitmq consume success", client.debugString())
	go client.listenCloseChan(ctx)
	return client, nil
}

func NewProducer(ctx context.Context, config URLConfig, mode int, exchangeName string, queueName string) (*Client, error) {
	conn, err := connect(ctx, config.GetDialURL())
	if err != nil {
		return nil, err
	}
	client := &Client{
		conn:         conn,
		connectURL:   config.GetDialURL(),
		mode:         mode,
		queueName:    queueName,
		exchangeName: exchangeName,
	}
	if err = client.initProducer(); err != nil {
		return nil, err
	}
	fmt.Println("rabbitmq producer success", client.debugString())
	go client.listenCloseChan(ctx)
	return client, nil
}

func (c *Client) listenCloseChan(ctx context.Context) {
	fmt.Println("rabbitmq listenCloseChan", c.debugString())
	closeSh := make(chan *amqp.Error)
	c.conn.NotifyClose(closeSh)
	err := <-closeSh
	if err != nil {
		fmt.Println("rabbitmq close", err, c.debugString())
	}

	//启动重连
	go c.reconnect(ctx)
}

func (c *Client) reconnect(ctx context.Context) {
	conn, err := connect(ctx, c.connectURL)
	if err != nil {
		fmt.Println("rabbitmq reconnect", err, c.debugString())
		time.Sleep(time.Second)
		go c.reconnect(ctx)
		return
	}

	fmt.Println("rabbitmq reconnect success", c.debugString())
	c.conn = conn
	if isConsumer {
		if err := c.initConsume(); err != nil {
			fmt.Println("rabbitmq reconnect success,initConsume fail", err, c.debugString())
		}
	} else {
		if err := c.initProducer(); err != nil {
			fmt.Println("rabbitmq reconnect success,initProducer fail", err, c.debugString())
		}
	}

	go c.listenCloseChan(ctx)
}

// initConsume
//
//	@Description: 初始化消费者
//	@receiver c
//	@return error
func (c *Client) initConsume() error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}
	switch c.mode {
	case AMQP_MODE_EXCHANGE:
		err = channel.ExchangeDeclare(
			c.exchangeName, // exchange名称
			"fanout",       // exchange类型
			true,           // 是否持久化
			false,          // 是否自动删除
			false,          // 是否内部
			false,          // 是否等待服务器确认
			nil,            // 其他参数
		)
		if err != nil {
			return fmt.Errorf("exchangeDeclare fail err:%v", err)
		}
		q, err := channel.QueueDeclare(
			"",    // 随机生成队列名称
			false, // 是否持久化
			true,  // 是否自动删除
			false, // 是否排他性
			false, // 是否等待服务器确认
			nil,   // 其他参数
		)
		if err != nil {
			return fmt.Errorf("queueDeclare fail err:%v", err)
		}
		//绑定
		err = channel.QueueBind(q.Name, "", c.exchangeName, false, nil)
		if err != nil {
			return fmt.Errorf("queueBind fail err:%v", err)
		}
		c.queueName = q.Name
	case AMQP_MODE_QUEUE:
		queue, err := channel.QueueDeclare(c.queueName, false, true, false, false, nil)
		if err != nil {
			return fmt.Errorf("queueDeclare fail err:%v", err)
		}
		c.queueName = queue.Name
	default:
		return fmt.Errorf("unhandled default case mode:%d", c.mode)
	}

	dataList, err := channel.Consume(c.queueName, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume fail err:%v", err)
	}
	go func() {
		defer fmt.Println("Consume close")
		for delivery := range dataList {
			if c.dataHandler == nil {
				fmt.Println("dataHandler is nil", c.debugString())
			} else {
				c.dataHandler(delivery.Body)
			}
		}
	}()
	return nil
}

// initProducer
//
//	@Description: 初始化生产者
//	@receiver c
//	@return error
func (c *Client) initProducer() error {
	var err error
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq channel err:%v", err)
	}
	//queue, err := c.channel.QueueDeclare(c.queueName, false, true, false, false, nil)
	//if err != nil {
	//	return fmt.Errorf("rabbitmq queue err:%v", err)
	//}
	//c.queueName = queue.Name
	return nil
}

func (c *Client) SendMessage(ctx context.Context, buf []byte) error {
	if c.channel == nil {
		return fmt.Errorf("rabbitmq channel is nil")
	}
	return c.channel.PublishWithContext(ctx, c.exchangeName, c.queueName, false, false, amqp.Publishing{
		Body: buf,
	})
}
