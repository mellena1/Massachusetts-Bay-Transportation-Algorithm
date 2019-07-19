import json
import matplotlib.pyplot as plt
import datetime

with open("../cubicSplineFunctions.json", 'r') as json_file:
    data = json.load(json_file)
    
    i = 0
    for edge in data:
        if i < 12:
            xs = []
            ys = []
            for pair in data[edge]['XYPairs']:    
                x = pair['X']
                y = pair['Y']
                if (x == 0 or  y == 0):
                    continue
                else:
                    ys.append(y)
                    x = x // 60
                    if x >= 24:
                        x = x - 24
                        now = datetime.datetime(2019, 7, 19, x)
                    else: 
                        now = datetime.datetime(2019, 7, 18, x)

                    xs.append(now.strftime('%-m/%-d %-H'))
            plt.plot(xs, ys, label=edge)
            i = i + 1
        else:
            break
        
    plt.gcf().autofmt_xdate()
    plt.legend()
    ax = plt.subplot(111)
    plt.xlabel('Date')
    plt.ylabel('Travel Time')
    ax.legend(loc='lower center', bbox_to_anchor=(0.5, -.27),
              ncol=3, fancybox=True, shadow=True)
    plt.show()
