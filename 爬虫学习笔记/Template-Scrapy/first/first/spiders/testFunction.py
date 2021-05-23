#!/usr/bin/env python3
# -*- coding: utf-8 -*-

__author__ = 'Chris'

from scrapy.spiders import Spider
from first.items import FirstItem

class FirstSpider(Spider):
    #使用scrapy crawl TEST_FUNCTION , 注意名字要匹配
    #这边三个参数的名字都不能变动，否则无法运行
    name = "TEST_FUNCTION"
    allowed_domains = ["sina.com.cn"]
    start_urls = ("https://www.sina.com.cn/",)

    def parse(self, response):
        """
        :param response: 该参数用于存放结果
        :return:
        """
        #使用items定义的参数，本质上是实例化类
        item = FirstItem()
        item["content"] = response.xpath("/html/head/title/text()").extract()
        #返回，不能使用return，否则程序会错误
        yield item