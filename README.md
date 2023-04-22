# Colorme CLI

`colorme` is an unofficial Color Me Shop command line tool.

## Installation

WIP

## Usage

### Log in

Log in to the Color Me Shop.

```
$ colorme login
```

Then, authenticate with your Color Me Shop account in the opened browser. If first time, authorize this tool to access your Color Me Shop account.

### Products

List products in the shop.

```
$ colorme product
=== Product 1
Name: マグカップ
Stocks: 20
Model Number: mug-001
Price: ¥1,760
Description: ベーシックなマグカップです。

View this product on the shop: https://test.shop-pro.jp/?pid=173944546

=== Product 2
Name: プレート
Stocks: 10
Model Number: plate-001
Price: ¥2,860
Description: 大きめのプレートです。

View this product on the shop: https://test.shop-pro.jp/?pid=168933092
```

### Orders

List orders in the shop.

```
$ colorme order
=== Order 154603080
Total Price: ¥1,920
Paid: false
Delivered: false
Canceled: false
View this order on the admin: https://admin.shop-pro.jp/?mode=shopcust_sales&sales_id=154603080

=== Order 154483308
Total Price: ¥3,120
Paid: true
Delivered: true
Canceled: false
View this order on the admin: https://admin.shop-pro.jp/?mode=shopcust_sales&sales_id=154483308
```

### Calling APIs

Make an authenticated request to [Color Me Shop API](https://developer.shop-pro.jp/docs/colorme-api).

Specify a path to the `/v1` API endpoint as the first argument.

```
$ colorme api /gift
{
  "gift": {
    "account_id": "PA00000000",
    "enabled": false,
    "noshi": {
      "enabled": false,
      "text_enabled": false,
      "types": [],
      "comment": null,
      "text_charge": null
    },
    "card": {
      "enabled": false,
      "text_enabled": null,
      "types": [],
      "comment": null
    },
    "wrapping": {
      "enabled": false,
      "types": [],
      "comment": null
    },
    "make_date": 1512918000,
    "update_date": 1512961758
  }
}
```
