	
text = "Hallo World Das ist ein Test Und so weiter Used"	
array = text.split(" ")
text_s = ""
value_s = array[len(array)-1]
for i in (range(len(array)-1)):
	text_s += array[i] +" "
print(text_s)
print(value_s)
