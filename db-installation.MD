sudo apt update
sudo apt install mariadb-server -y
sudo /etc/init.d/mysql start
sudo /etc/init.d/mysql status
sudo mysql_secure_intallation
sudo mysql -u root -p

create database inventory;

use inventory;

create table products(
id int NOT NULL AUTO_INCREMENT ,
name varchar(255) NOT NULL,
quantity int,
price float(10,7),
PRIMARY KEY(id)
);

insert into products values(1 , "chair", 100, 200.00);
insert into products values(2 , "table", 300, 150.00);

select * from products;

go mod init example.com/my-inventory
go get github.com/gorilla/mux
go get github.com/go-sql-driver/mysql