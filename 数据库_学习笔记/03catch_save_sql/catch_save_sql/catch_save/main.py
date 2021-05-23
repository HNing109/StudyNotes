#!/usr/bin/env python3
# -*- coding: utf-8 -*-

__author__ = 'Chris'

from scrapy.cmdline import execute
import os
import sys

sys.path.append(os.path.dirname(os.path.abspath(__file__)))
#显示执行的提示详情
# execute(["scrapy", "crawl", "CATCH"])
#只显示结果
execute(["scrapy", "crawl", "CATCH", "--nolog"])
