Hello World, I'm totally new here
Purpose of this project is to learning hexagonal architecture and Go syntax

bash
docker-compose up -d = to run our Mongodb

User routes:
post http://localhost:8080/register = sign up with username admin , to have authority to every content
json
{"username":"admin" ,"password":"admin","email":"asd@hotmail.com"}

post http://localhost:8080/login = login with username and password, returns accesstoken and refreshtoken. use accesstoken as Bearer token to get authorization.
json
{"username":"admin","password":"admin"}

get http://localhost:8080/user/getuser = get user by jwt token that contains userid

post http://localhost:8080/user/update = update user informations
json
{"email":"asd@hotmail.com","firstname":"firstname","lastname":"lastname","address":"address"}

post http://localhost:8080/user/resetpassword = gonna send you a new genrated password through email service
json
{"email":"yourregisteredemail@gmail.com"}

post http://localhost:8080/user/changepassword
{"oldpassword":"oldpass","newpassword":"newpass"}

get http://localhost:8080/user/cart/getcart = get users cart , show products in the cart, using userid in token

post http://localhost:8080/user/addtocart = this gonna add product to cart
json
{"productid":"66c96f7f600fb3ed9788a3f7","productname":"productname","image":["1.jpg"],"details":"details","amount":1,"priceeach":1500,"priceid":"price_1PrCWZHzxBcAUB6C0H29ahiE"}

post http://localhost:8080/user/cart/increase = this gonna increase product in cart by 1, imagine as a + button at cart page

post http://localhost:8080/user/cart/decrease

post http://localhost:8080/user/cart/deleteproduct = use productid to delete product from user's cart

post http://localhost:8080/user/getorder = get user orders by userid that claims from token

Product route:
get http://localhost:8080/product/all = get all product

For only admin
post http://localhost:8080/private/product/:productid = searching product by ID

post http://localhost:8080/private/product/new = create a new product, this gonna create productid and priceid at stripe, store product informations on stripe that might use in the future
json
{"productname":"productname","details":"details","stock":1,"category":"categroy","image":["1","2"],"priceeach":22}

post http://localhost:8080/private/product/update = updates product informations, this gonna update details on stripe as well
json
{"productname":"rod","details":"newdetails","stock":1,"category":"weapon","image":["1","2"],"priceeach":22,"stripeproductid":"prod_Qibr1S3soTIx2r","productid":"66c9534484e386d058a89c57"}

post http://localhost:8080/private/product/delete = delete product in product collection
