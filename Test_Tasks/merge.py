# -*- coding: utf-8 -*-

def merge_sort(array):  

    def merge(a, left, m, right):
        # left, m, right - индексы
        if left >= right: 
            return None
        if m < left or right < m: 
            return None

        t = left
        for j in range(m + 1, right + 1): # подгруппа 2
            for i in range(t, j): # цикл подгруппы 1
                if a[j] < a[i]:
                    r = a[j]
                    # итерационно переставляем элементы, чтобы упорядочить
                    for k in range(j, i, -1):
                        a[k] = a[k - 1]
                    a[i] = r
                    t = i # проджолжение вставки в группе 1
                    break # к следующему узлу из подгруппы 2
                
    if len(array) < 2:
     return None

    k = 1 # шаг
    while k < len(array):
        g = 0 # 
        while g < len(array): # группы
            z = g + k + k - 1 # последний эл-т группы
            r = z if z < len(array) else len(array) - 1 # последняя группа
            merge(array, g, g + k - 1, r) # слияние
            g += 2*k
        k*=2

print("*** Merge sort ***")
numbers = raw_input("Enter numbers splitted by space:").split()
numbers = [int(el) for el in numbers]

merge_sort(numbers)

print(numbers)
