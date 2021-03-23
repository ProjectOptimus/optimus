import json

x: list = json.loads('{"a": 1, "b": 2}')

y = x['a'] + x['b']
