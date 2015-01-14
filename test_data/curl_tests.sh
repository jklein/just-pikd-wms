#curl tests to be run against go server to test various endpoints (later many of these should become unit tests)

#test data reset
curl -4 -i -X POST http://localhost:3000/reset

#test PUTing receiving location - can change the supplier_shipment_id as well
curl -4 -i -X PUT http://localhost:3000/locations/receiving -d '{
    "receiving_location_id": "204-178900284-5",
    "receiving_location_type": "Pallet Receiving",
    "temperature_zone": "dry",
    "supplier_shipment_id": 14
  }'

#should return a big list of shipments
curl -4 -i http://localhost:3000/suppliers/shipments

#should return a single id
curl -4 -i http://localhost:3000/suppliers/shipments?shipment_id=397787316

#should return a single id
curl -4 -i http://localhost:3000/suppliers/shipments?shipment_id=397787316&stocking_purchase_order_id=2

#should return 404
curl -4 -i http://localhost:3000/suppliers/shipments?shipment_id=397787316&stocking_purchase_order_id=1

#should return 200
curl -4 -i -X PUT http://localhost:3000/suppliers/shipments/4 -d '{
    "supplier_shipment_id": 4,
    "shipment_id": "397787316",
    "stocking_purchase_order_id": 0,
    "supplier_id": 2,
    "promised_delivery": null,
    "actual_delivery": "2015-01-17T00:00:00Z"
  }'

#should return 400 bad request
curl -4 -i -X PUT http://localhost:3000/suppliers/shipments/5 -d '{
    "supplier_shipment_id": 4,
    "shipment_id": "397787316",
    "stocking_purchase_order_id": 0,
    "supplier_id": 2,
    "promised_delivery": null,
    "actual_delivery": "2015-01-17T00:00:00Z"
  }'
