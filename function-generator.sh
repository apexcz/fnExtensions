#!/bin/bash
echo "Generating multiple functions."

cd ..
for i in {1..10}
do
  fn init --runtime go droplet_check$i
  cd droplet_check$i
  export FN_REGISTRY=bigty
  fn deploy --app check
  echo "deployed round $i"
  cd ..
done

echo "Finished deployment"
