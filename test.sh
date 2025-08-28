#curl -X POST \
#  -H "Content-Type: application/json" \
#  -H "Authorization: Bearer 8645b0c7b819fd54f7d95c6b80445681" \
#  -d '{"query": "{positions(where: {owner: \"0xE36EF7f03Be54AEAeCdbF5A1a69a21620fb09686\", liquidity_gt: 0}){id,tickLower{tickIdx},tickUpper{tickIdx},pool{tick,token0{symbol},token1{symbol}}}}", "operationName": "Subgraphs", "variables": {}}' \
#  https://gateway.thegraph.com/api/subgraphs/id/GENunSHWLBXm59mBSgPzQ8metBEp9YDfdqwFr91Av1UM

curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 8645b0c7b819fd54f7d95c6b80445681" \
  -d '{"query": "{positions(where: {owner: \"0xE36EF7f03Be54AEAeCdbF5A1a69a21620fb09686\", liquidity_gt: 0}){id,tickLower{tickIdx},tickUpper{tickIdx},pool{tick,token0{symbol},token1{symbol}}}}", "operationName": "Subgraphs", "variables": {}}' \
  https://gateway.thegraph.com/api/subgraphs/id/GENunSHWLBXm59mBSgPzQ8metBEp9YDfdqwFr91Av1UM

