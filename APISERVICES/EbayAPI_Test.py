
# coding=utf-8

import time
from ebaysdk.exception import ConnectionError
from ebaysdk.finding import Connection
from collections import Counter
import numpy as np
from scipy.interpolate import spline
import statistics
import numpy as np
from sklearn.cluster import MeanShift, estimate_bandwidth
from sklearn.cluster import KMeans


def Ebay_Calc(Begriff,value_s):

    eintraege = {}
    
    eintraege_Liste = []
    start_time = time.time()
    
    prices = []
    prices_gerundet = []
    try:
        api = Connection(appid='BrusLan-private-PRD-d3391a9b0-f8060535', config_file=None,siteid='EBAY-DE')
        print("--- %s seconds ---" % (time.time() - start_time))
        response = api.execute('findCompletedItems', {'keywords': '{}'.format(Begriff),'itemFilter':[{'name': 'Condition','value': value_s},{'name': 'TopRatedSellerOnly', 'value': 'false'},{"name":"SoldItemsOnly","value":"true"},{"name":"Currency","value":"EUR"}],"paginationInput":{"entriesPerPage":"100","pageNumber":"1"}})
        print("--- %s seconds ---" % (time.time() - start_time))
        
        eintraege = response.dict()
        
        try:
            eintraege_Liste  = eintraege["searchResult"]["item"]
        except:
            return 0
      
        # print(type(eintraege_Liste))
        # if (len(eintraege_Liste)==100):
        # 	print("sind 100")
        # 	try:
	       #  	response = api.execute('findCompletedItems', {'keywords': '{}'.format(Begriff),'itemFilter':[{'name': 'Condition','value': value_s},{'name': 'TopRatedSellerOnly', 'value': 'false'},{"name":"SoldItemsOnly","value":"true"},{"name":"Currency","value":"EUR"}],"paginationInput":{"entriesPerPage":"100","pageNumber":"2"}})
	       #  	eintraege = response.dict()
	       #  	eintraege_Liste.append(eintraege["searchResult"]["item"])

	       #  except:
	       #  	print("hat halt nicht geklappt")
       
        
        #print(len(eintraege["searchResult"]["item"]))
        print(len(eintraege_Liste))


    except ConnectionError as e:
        return print("Error in Connections",e)



    for Eintrag in eintraege_Liste:
        #print(Eintrag)
        if (Eintrag["condition"]["conditionId"] != "7000"):
            prices.append(int(float(Eintrag["sellingStatus"]["currentPrice"]["value"])))
        else:
            print("defekt")
        #print("\n")




    for price in prices:
        prices_gerundet.append(round(price/5)*5)


    unique_price = Counter(prices)
    unique_price_gerundet = Counter(prices_gerundet)
    try:
        clusterpreis = clustering(prices)
    except: 
        return 0
    #plt.scatter(prices,prices)
    #sys.exit(0)
    # plt.bar(unique_price_gerundet.keys(), unique_price_gerundet.values())
    #plt.show()
    # t = np.arange(0,max(prices))

    #Plotting hours of feeds
    # y = np.zeros(max(prices))
    # for hour in t:
    #   if hour in unique_price.keys():
    #       y[hour] = unique_price[hour]
    #   else:
    #       np.append(y, 0)
    # xnew = np.linspace(t.min(),t.max(),300)
    # power_smooth = spline(t,y[1:],xnew)
    print("--- %s seconds ---" % (time.time() - start_time))
    print(clusterpreis)
    return clusterpreis

    # plt.plot(xnew,power_smooth,"k",t,y[1:],"bo")
    # # eintraege = pickle.load(pickle_file)
    # # print(eintraege)
    # plt.show()
    #pickle.dump(eintraege,pickle_file)
def clustering(x):



    X = np.array(list(zip(x,np.zeros(len(x)))), dtype=np.int)
    bandwidth = estimate_bandwidth(X, quantile=0.1)
    print(bandwidth)
    ms = MeanShift(bandwidth=bandwidth, bin_seeding=True)
    ms.fit(X)
    labels = ms.labels_
    cluster_centers = ms.cluster_centers_

    labels_unique = np.unique(labels)
    n_clusters_ = len(labels_unique)

    clusterList = []

    for k in range(n_clusters_):
        my_members = labels == k
        print ("cluster {0}: {1}".format(k, X[my_members, 0]))
        print(sum(X[my_members, 0]) / float(len(X[my_members, 0])))
        print(statistics.median(X[my_members, 0]))
        clusterList.append(statistics.median(X[my_members, 0]))
     
    print(clusterList)
    myList = clusterList

    print(reject_outliers(np.array(myList)))
    returnPreis  = reject_outliers(np.array(myList))[0]


    return returnPreis

     


def reject_outliers(data, m=1):
    return data[abs(data - np.mean(data)) < m * np.std(data)]


if __name__ == "__main__":
	print(Ebay_Calc("Entsperrt Apple iPad Pro 32GB Wi-Fi + Cellular 4G LTE 9.7 Space grau","Used"))

    

