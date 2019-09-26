## Another data visualisation app for strava

Just because I'm curious and I want to... I've decided to show how the KJ that I've generated, can be transformed into power, while I know it might not be
entirely perfect, the idea would be a nice approximation

See for better details on how the authentication is done http://developers.strava.com/docs/authentication/

For now, pick and axe solution:

```
source production.env
export APP_SCOPES="activity:read_all,profile:read_all,activity:write,read_all,read"
#this returns an URL, go there, extract code and put it as AC
curl -X GET "http://www.strava.com/oauth/authorize?client_id=$CID&response_type=code&redirect_uri=http://localhost/exchange_token&approval_prompt=force&scope=$APP_SCOPES"
export AC="RESULT OF REDIRECTION's URL"

#this returns json, get the token from here
curl -X POST https://www.strava.com/oauth/token  -F client_id=$CID -F client_secret=$CS  -F code=$AC  -F grant_type=authorization_code
echo export STRAVA_TOKEN="TOKEN_GIVEN" > production.env
source production.env 
curl  -D /dev/stderr -X GET  -H "Authorization: Bearer $STRAVA_TOKEN"  'https://www.strava.com/api/v3/athlete/activities'  | tee activities.json | jq " map(. + {kwh: ( (.kilojoules // 0) / 3.6)} ) "
```

## What to represent?

* Ideally I want to show how many light bulbs (90w or some fancy LED) can be powered in an hour

### further reading

* https://m.wikihow.com/Calculate-Kilowatts-Used-by-Light-Bulbs
* https://insights.regencylighting.com/kw-vs-kwh-how-much-energy-is-my-lighting-using
* http://energyusecalculator.com/electricity_cfllightbulb.htm
* https://www.arcadiapower.com/energy-101/energy-bills/how-to-visualize-one-kwh/
