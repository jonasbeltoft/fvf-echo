# Writes 100.000 entries to test_data.csv
# Each entry has 3 columns: name, data, link, fill

# name: random string of length 10
# data: random array of floats of length 2 - 5
# link: random string of length 20
# fill: random string of length 127

import random
import string
import csv

def random_string(length):
	return ''.join(random.choice(string.ascii_lowercase) for i in range(length))

# Open file or create it if it doesn't exist
with open('data/test_data.csv', 'w', newline='') as csvfile:
	writer = csv.writer(csvfile)
	writer.writerow(['name', 'data', 'link', 'fill'])
	for i in range(10000):
		writer.writerow([random_string(10), [random.random() for i in range(random.randint(2, 5))], random_string(20), random_string(127)])