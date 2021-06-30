# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class CatchSaveItem(scrapy.Item):
    # define the fields for your item here like:
    # name = scrapy.Field()
    print("------------RUN item.py-------------")

    title = scrapy.Field();
    url = scrapy.Field();
    comment = scrapy.Field();
    price = scrapy.Field();

