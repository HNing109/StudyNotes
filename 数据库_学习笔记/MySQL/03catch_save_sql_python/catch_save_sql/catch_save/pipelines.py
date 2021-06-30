# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html

import pymysql

server = "192.168.1.103"
user = "chris"
password = "1230re0321re"
database = "dangdang"

class CatchSavePipeline(object):
    def process_item(self, item, spider):
        # 打开数据库连接
        db = pymysql.connect(server, user, password,database)
        # 使用cursor()方法获取操作游标
        cursor = db.cursor()
        for i in range(0, len(item["title"])):
            print(i)
            title = item["title"][i]
            price = item["price"][i]
            comment = item["comment"][i]
            url = item["url"][i]
            sql = "insert into bookinfo(title, price, comment, url) \
                   values ('%s',  '%s',  '%s',  '%s')" % \
                  (title, price, comment, url)
            print(title + "  ---  " + price + "  ---  " + comment + url + "\n")
            # 执行sql语句
            cursor.execute(sql)
            # 提交到数据库执行
            db.commit()
        # 关闭数据库
        db.close()
        return item

    # def __init__(self):
    #     self.df = open("./result/spider.txt","a")

    # def process_item(self, item, spider):
    #     print(item["title"])
    #     print(item["url"])
    #     print(item["comment"])
    #     print("-----------")
    #     self.df.write(item["title"][0] + "\n" + item["url"][0] + "\n" + item["comment"][0] + "\n------------\n")
    #     return item


    # 这是默认最后执行的方法
    # def close_spider(self):
    #     self.df.close()