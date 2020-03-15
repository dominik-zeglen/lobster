import matplotlib.pyplot as plt
import numpy as np
import json

chunkSize = 1024
rate = 44100

with open('../output.json') as json_file:
    data = json.load(json_file)

plt.plot(data[(chunkSize * 10 - 200):(chunkSize * 10 + 200)])
plt.savefig("plot.svg", format="svg")