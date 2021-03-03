# 数字の配列が渡されます
# その配列から「二乗した値が2で割り切れるもの」を抽出し文字列の配列として返す関数
def f(l):
    return list(reversed([f"{i}" for i in l if (i ** 2) % 2 == 0]))

print(f([1, 2, 3, 4]))
