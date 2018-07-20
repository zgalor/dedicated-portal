create database customers;
create table customers (
  id             text not null unique primary key,
  name           text not null,
);
create table owned_clusters (
  customer_id  text not null references customers (id),
  cluster_id   text not null unique
);
