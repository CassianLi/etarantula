package app

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
	"tarantula-v2/models"
	rabbitmq "tarantula-v2/rabbit"
	"tarantula-v2/service"
	"time"
)

// 消息处理函数
func handler(msg string) {
	log.Println("收到MQ消息：", msg)

	start := time.Now()
	var category models.CategoryInfoRequest
	// 将 JSON 字符串转换为 Go 类型的实例
	err := json.Unmarshal([]byte(msg), &category)
	if err != nil || category.Country == "" || category.ProductNo == "" {
		fmt.Println("Error unmarshalling message, err: ", err, category)
		return
	}
	// 转换为大写
	category.Country = strings.ToUpper(category.Country)

	s := service.NewCategoryService(category)
	if s == nil {
		log.Println("Create service failed")
		return

	}

	info, err := s.GetCategoryInfo()
	// 调试时，打印回传信息
	if viper.GetBool("debug") {
		// categoryInfo 转json string
		categoryInfoJson, err := json.Marshal(info)
		if err != nil {
			log.Println("即将回传截图消息MQ, categoryInfo: ", info)
		} else {
			log.Println("即将回传截图消息MQ, categoryInfo(json): ", string(categoryInfoJson))
		}
	}
	end := time.Now()
	log.Println("获取品类截图信息总耗时：", end.Sub(start))

	start = time.Now()
	// publish category info to MQ
	err = publishInfo(info)
	if err != nil {
		log.Println("回传截图信息，失败, err: ", err)
		return
	}
	end = time.Now()
	log.Println("回传截图信息总耗时:", end.Sub(start))

}

// publishInfo 发布回传信息到MQ
func publishInfo(info models.CategoryInfo) error {
	// type 转json string
	infoJson, err := json.Marshal(info)
	if err != nil {
		fmt.Println("Error marshalling message, err: ", err)
		return err
	}

	mq, err := rabbitmq.NewRabbitMQ(viper.GetString("pub-mq.url"), viper.GetString("pub-mq.exchange"), viper.GetString("pub-mq.exchange-type"))
	if err != nil {
		fmt.Println("创建消息回传MQ链接失败，Error: ", err)
		return err
	}
	defer mq.Close()

	err = mq.Publish(viper.GetString("pub-mq.queue"), infoJson)
	if err != nil {
		fmt.Println("Error publishing message, err: ", err)
		return err
	}
	return nil
}

// Consuming 启动消费者
func Consuming(url, exchange, exchangeType, queue string) {
	mq, err := rabbitmq.NewRabbitMQ(url, exchange, exchangeType)
	defer mq.Close()

	if err != nil {
		fmt.Println("创建MQ链接失败，Error: ", err)
		return
	}

	fmt.Println("创建MQ链接成功，开始消费...")
	err = mq.Consume(queue, false, handler)

	if err != nil {
		fmt.Println("消费MQ消息失败，Error: ", err)
		return
	}
}
