from bottle import route, run, template, response, request, get
from json import dumps
import bottle
import json
import EbayAPI_Test


def testfuc(text_s, value_s):

	return str(EbayAPI_Test.Ebay_Calc(text_s, value_s))

@get('/test/<text>')
def test(text):

	text = text.replace('"', "")
	print(text)
	array = text.split(" ")
	text_s = ""
	value_s = array[len(array)-1]
	for i in (range(len(array)-1)):
		text_s += array[i] +" "
	print(text_s)
	print(value_s)
	return testfuc(text_s, value_s)

@get('/hallo')
def test2():
	
	return "hallo"




run(host = "0.0.0.0", port = 8080, debug=True)
 



