
#!/bin/bash

echo "########### Loading data to Quotes Mongo DB Collection ###########"
mongoimport --jsonArray --db micromind --collection quotes --file /tmp/data/quotes.json


#!/bin/bash

echo "########### Loading data to Questions Mongo DB Collection ###########"
mongoimport --jsonArray --db micromind --collection questions --file /tmp/data/questions.json