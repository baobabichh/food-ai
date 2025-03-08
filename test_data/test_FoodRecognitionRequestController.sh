curl -X GET http://localhost:5555/createRequest \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "uuid=065cee3f-2d88-4983-b753-d754677eb2ba" \
     --data-urlencode "img_base64@1.txt" \
     -d "img_type=jpg" \
     -i \


curl -X GET http://localhost:5555/checkRequestStatus \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "uuid=065cee3f-2d88-4983-b753-d754677eb2ba" \
     -d "request_id=4" \
     -i \

curl -X GET http://localhost:5555/getRequest \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "uuid=065cee3f-2d88-4983-b753-d754677eb2ba" \
     -d "request_id=4" \
     -i \