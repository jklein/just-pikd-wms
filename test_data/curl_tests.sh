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

#should succeed
curl -4 -i -X PATCH http://localhost:3000/spos/1 -d '{
    "spo_id": 1,
    "spo_date_arrived": "2015-01-14T00:00:00Z"
}'

#should 400
curl -4 -i -X PATCH http://localhost:3000/spos/2 -d '{
    "spo_id": 1,
    "spo_date_arrived": "2015-01-14T00:00:00Z"
}'

#should 404
curl -4 -i -X PATCH http://localhost:3000/spos/2000 -d '{
    "spo_id": 2000,
    "spo_date_arrived": "2015-01-14T00:00:00Z"
}'

#should succeed
curl -4 -i -X PATCH http://localhost:3000/spos/1 -d '{
    "spo_id": 1,
    "products": [
      {
        "spop_id": 109,
        "spop_confirmed_qty": 2,
        "spop_received_qty": 2
      }
    ]
}'

#should succeed and update both the SPO and two products
curl -4 -i -X PATCH http://localhost:3000/spos/1 -d '{
    "spo_id": 1,
    "spo_date_arrived": "2015-01-16T00:00:00Z",
    "products": [
      {
        "spop_id": 109,
        "spop_confirmed_qty": 3,
        "spop_received_qty": 3
      },
      {
        "spop_id": 111,
        "spop_confirmed_qty": 5,
        "spop_received_qty": 5
      }
    ]
}'

#should 404
curl -4 -i -X PATCH http://localhost:3000/spos/1 -d '{
    "spo_id": 1,
    "spo_date_arrived": "2015-01-14T00:00:00Z",
    "products": [
      {
        "spop_id": 123123123,
        "spop_confirmed_qty": 3,
        "spop_received_qty": 3
      },
      {
        "spop_id": 111,
        "spop_confirmed_qty": 2,
        "spop_received_qty": 2
      }
    ]
}'

#leaving the ID out of the embedded document entirely should also 404
curl -4 -i -X PATCH http://localhost:3000/spos/1 -d '{
    "spo_id": 1,
    "spo_date_arrived": "2015-01-14T00:00:00Z",
    "products": [
      {
        "spop_confirmed_qty": 3,
        "spop_received_qty": 3
      }
    ]
}'

#should get a whole lot of results
curl -4 -i -X GET 'http://localhost:3000/spos?supplier_id=1'

#should get a 404
curl -4 -i -X GET 'http://localhost:3000/spos?supplier_id=2'

#should get one spo
curl -4 -i -X GET 'http://localhost:3000/spos?supplier_id=1&shipment_code=803207203'