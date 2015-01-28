#curl tests to be run against go server to test various endpoints (later many of these should become unit tests)

#test data reset
curl -4 -i -X POST http://localhost:3000/reset

#test PUTing receiving location - can change the supplier_shipment_id as well
curl -4 -i -X PATCH http://localhost:3000/locations/receiving/204-178900284-5 -d '{
    "rcl_id": "204-178900284-5",
    "rcl_shi_shipment_code": 14
  }'

#should return a big list of shipments
curl -4 -i http://localhost:3000/suppliers/shipments

#should return a single id
curl -4 -i 'http://localhost:3000/suppliers/shipments?shipment_code=777152188'

#should return a single id
curl -4 -i 'http://localhost:3000/suppliers/shipments?shipment_code=777152188&spo_id=3'

#should return 404
curl -4 -i 'http://localhost:3000/suppliers/shipments?shipment_code=777152188&spo_id=1'

#should return 200
curl -4 -i -X PATCH http://localhost:3000/suppliers/shipments/4 -d '{
    "shi_id": 4,
    "shi_actual_delivery": "2015-01-14T00:00:00Z"
  }'

#should return 400 bad request
curl -4 -i -X PATCH http://localhost:3000/suppliers/shipments/5 -d '{
    "shi_id": 4,
    "shi_actual_delivery": "2015-01-17T00:00:00Z"
}'

#should 404
curl -4 -i -X PATCH http://localhost:3000/suppliers/shipments/10000 -d '{
    "shi_id": 10000,
    "shi_actual_delivery": "2015-01-17T00:00:00Z"
}'

#should get a whole lot of results
curl -4 -i -X GET 'http://localhost:3000/spos?supplier_id=1'

#should get a 404
curl -4 -i -X GET 'http://localhost:3000/spos?supplier_id=2'

#should get one spo
curl -4 -i -X GET 'http://localhost:3000/spos?supplier_id=1&shipment_code=803207203'