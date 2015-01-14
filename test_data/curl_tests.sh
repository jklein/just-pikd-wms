#curl tests to be run against go server to test various endpoints

#test data reset
curl -4 -i -X POST http://localhost:3000/reset

#test PUTing receiving location - can change the supplier_shipment_id as well
curl -4 -i -X PUT localhost:3000/locations/receiving -d '{
    "receiving_location_id": "204-178900284-5",
    "receiving_location_type": "Pallet Receiving",
    "temperature_zone": "dry",
    "supplier_shipment_id": 14
  }'