#!/usr/bin/env python3
# -*- coding: utf-8 -*-

__author__ = 'Chris'

from scrapy.spiders import Spider
from scrapy.http import Request
from scrapy.spiders import Rule, CrawlSpider
from scrapy.linkextractors import LinkExtractor
from catch_save.items import CatchSaveItem


class Catch(Spider):
    name = "CATCH"
    allowed_domains = ["dangdang.com"]
    Homepage = 'http://www.dangdang.com/'

    def start_requests(self):
        """
        第一次爬取，会使用start_requests的模拟浏览器，后面的会使用setting.py的User-Agent
        :return:
        """
        print("-------------RUN start_requests---------------")
        headers = {"User-Agent":'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36'}
        yield Request(self.Homepage, headers=headers)

    # rules = (
    #     Rule(LinkExtractor(allow = r'article'), callback='parse', follow=True),
    # )

    def parse(self, response):
        print("-------------------RUN parse------------------")
        item = CatchSaveItem()
        item["title"] = response.xpath("//a[@class='pic']/@title").extract()
        item["price"] = response.xpath("//span[@class='search_now_price']/text()").extract()
        # @href : 提取href标签的值
        item["url"] = response.xpath("//a[@class='pic']/@href").extract()
        item["comment"] = response.xpath("//a[@class='search_comment_num']/text()").extract()

        yield item

        for i in range(1,40):
            url = "http://category.dangdang.com/pg"+ str(i) + "-cp01.52.04.02.00.00.html"
            yield Request(url, callback=self.parse)

