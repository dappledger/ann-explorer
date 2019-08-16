# Prepare

Download go-ethereum <br/>
Install mongo if you want store mass historical data or use another machine to store data


# Build

make

Look `.gitlab-ci.yml` if you make failed

# Note

Modify `conf/app.conf` ,ensure `api_addr` serve is running <br/>
Add `chain_id` in `conf/app.conf` if your system support multiple blockchains (such as `civilization`) <br/>
Add mongo config in `conf/app.conf`,if you want use mongo

# Run

./block-browser

# Example

Exec shell `sed -i "/^api_addr/d" conf/app.conf ;echo 'api_addr = "10.253.4.248:6657"' >> conf/app.conf ;./block-browser`
