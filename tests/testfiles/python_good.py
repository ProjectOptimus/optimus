"""
this is a module docstring
"""

import json

x: dict = json.loads('{"a": 1, "b": 2}')

y = x["a"] + x["b"]
