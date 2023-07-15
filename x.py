import base64, json
import quickjs, httpx, math

def get_token():
    b64 = httpx.get("https://sdk.vercel.ai/openai.jpeg").text
    data = json.loads(base64.b64decode(b64))

    script = """
      var globalThis = {{data: "sentinel"}};
      ({script})({key})
    """.format(script=data["c"], key=data["a"])
    
    context = quickjs.Context()
    token_data = json.loads(context.eval(script).json())

    token = {
      "r": token_data,
      "t": data["t"]
    }

    token_string = json.dumps(token, separators=(',', ':')).encode()

    print(data)
    print(token_string)
    print(data["a"] % math.floor(data["a"] * math.exp(1)))

    # {'r': [3.299471618955463, [], None], 't': 'eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..AoFuobOe_3tpP3LJ.YUpjCVBYVCWyJ37lAoBXY_OvHsbIMD68K06hsQ2el995XWkH8w3lQ-jIFsysVY5CnowmZ45afHwwcxL0ytRjrp12UPW97HcQw7BN6IXQi0-ypBIuLxB_yL8qsoqttYA.SJ6T6EJ5NvrF96LqCd6fxQ'
    return base64.b64encode(token_string).decode()

#print(get_token())
get_token()