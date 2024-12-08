CREATE database my_db;

USE my_db;
DROP TABLE Users;
CREATE TABLE Users 
(ID int PRIMARY KEY AUTO_INCREMENT, 
Name VARCHAR(100), 
Email VARCHAR(100), 
Password VARCHAR(100),
MemberTier VARCHAR(10),
Bookings int);

SELECT * FROM Users;

DROP TABLE Cars;
CREATE TABLE Cars
(ID int PRIMARY KEY AUTO_INCREMENT,
Name VARCHAR (100),
Description VARCHAR (5000),
ImageLink VARCHAR (1000),
PriceHour int,
MemberTier VARCHAR (100));

INSERT INTO Cars (Name, Description, ImageLink, PriceHour, MemberTier) 
VALUES ("Hyundai IONIQ Electric", "Permanent magnet synchronous, single speed reduction gear.", "https://media.ed.edmunds-media.com/hyundai/ioniq-electric/2020/oem/2020_hyundai_ioniq-electric_4dr-hatchback_limited_fq_oem_1_1600.jpg",
 10, "Basic"),
 ("BYD E6", "AC permanent magnet synchronous motor", "https://stimg.cardekho.com/images/carexteriorimages/930x620/BYD/E6/8675/1676289700142/front-left-side-47.jpg",
 15, "Basic"),
 ("BYD ATTO 3", "Permanent magnet synchronous motor", "https://www.bydauto.co.nz/storage/uploads/112972f5-aee9-4ddc-b1a7-19fb846a9818/model-page-banner-atto-3-2160x1185-q1-2024-02.jpg",
 20, "Premium"),
("XPENG G6", "Ultra-low energy consumption and a reliable driving range of 570 WLTP, with a charging capacity of up to 280 kW", "https://s-cdn.xpeng.com/xpwebsite/prod/2024-04-03/884d01a2936e4e2283c690797a827cd7.jpeg",
25, "Premium"),
("TESLA Model Y", "Dual motor all-wheel,  with up to 533 km (WLTP) of range on a single charge", "https://media.ed.edmunds-media.com/tesla/model-y/2024/oem/2024_tesla_model-y_4dr-suv_performance_fq_oem_1_1600.jpg",
40, "VIP");

SELECT * FROM Cars;

DROP TABLE CarsBooking;
CREATE TABLE CarsBooking
(ID int PRIMARY KEY AUTO_INCREMENT,
UserID int,
CarName VARCHAR (100),
CarID int,
StartDate DATETIME,
EndDate DATETIME,
TotalHours float,
TotalCost float,
Status VARCHAR(100));

SELECT * FROM CarsBooking;

DROP TABLE Billing;
CREATE TABLE Billing
(ID int PRIMARY KEY AUTO_INCREMENT,
UserID int,
CarID int,
BookingID int,
CarName VARCHAR(100),
StartDate DATETIME,
EndDate DATETIME,
PriceHour int,
TotalCost float,
UserCard VARCHAR(100),
Status VARCHAR(100));

SELECT * FROM Billing;