# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html


class FirstPipeline(object):
    def process_item(self, item, spider):
        #item：item.py中的容器，输出对应的标签：content
        print(item["content"])
        return item
