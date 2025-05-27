CREATE TABLE "Venues" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "image_links" varchar[],
  "name" varchar(255) NOT NULL,
  "type" varchar(255),
  "description" varchar(255),
  "location" varchar(255) NOT NULL,
  "dimension" varchar(255),
  "capacity" integer,
  "facilities" varchar[],
  "has_accomodation" boolean DEFAULT (false),
  "room_type" varchar(255),
  "no_of_rooms" integer,
  "sleeps" varchar(255),
  "bed_type" varchar(255),
  "rent" integer,
  "owned_by" uuid,
  "is_available" boolean DEFAULT (false),
  "opens_at" timestamp,
  "closes_at" timestamp,
  "rental_days" varchar(255),
  "booking_price" integer,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "Profiles" (
  "id" uuid PRIMARY KEY,
  "bio" varchar(255),
  "phone_no" integer,
  "country" varchar(255),
  "address" varchar(255),
  "experience" integer,
  "field" varchar(255),
  "business_name" varchar(255),
  "roles" string NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "Practitioners" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "name" varchar(255) NOT NULL,
  "description" varchar(255) NOT NULL,
  "image_link" varchar(255),
  "is_available" boolean DEFAULT (true),
  "created_by" uuid,
  "opens_at" timestamp,
  "closes_at" timestamp,
  "working_days" varchar(255),
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "Events" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "venue_id" uuid,
  "image_links" varchar[],
  "name" varchar(255),
  "theme" varchar(255),
  "description" varchar(255),
  "audience" varchar(255),
  "activities" varchar[],
  "created_by" uuid,
  "start_time" timestamp,
  "start_date" date,
  "end_date" date,
  "total_particpant" integer,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "Favourites" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "event_id" uuid,
  "added_by" uuid,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "Users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "email" varchar(255) UNIQUE NOT NULL,
  "password" varchar(255),
  "firstname" varchar(255),
  "lastname" varchar(255),
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "BookedVenues" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "type" varchar(255),
  "venue_id" uuid,
  "booked_for" datetime,
  "booked_by" uuid,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "BookedPractitioners" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "type" varchar(255),
  "service_id" uuid,
  "booked_for" datetime,
  "booked_by" uuid,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "Purchases" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "event_id" uuid,
  "venue_id" uuid,
  "service_id" uuid,
  "purchased_by" uuid,
  "created_at" timestamp DEFAULT (now())
);

ALTER TABLE "Profiles" ADD FOREIGN KEY ("id") REFERENCES "Venues" ("owned_by");

ALTER TABLE "Profiles" ADD FOREIGN KEY ("id") REFERENCES "Events" ("created_by");

ALTER TABLE "Venues" ADD FOREIGN KEY ("id") REFERENCES "Events" ("venue_id");

CREATE TABLE "Profiles_Users" (
  "Profiles_id" uuid,
  "Users_id" uuid,
  PRIMARY KEY ("Profiles_id", "Users_id")
);

ALTER TABLE "Profiles_Users" ADD FOREIGN KEY ("Profiles_id") REFERENCES "Profiles" ("id");

ALTER TABLE "Profiles_Users" ADD FOREIGN KEY ("Users_id") REFERENCES "Users" ("id");


ALTER TABLE "Events" ADD FOREIGN KEY ("id") REFERENCES "Favourites" ("event_id");

ALTER TABLE "Venues" ADD FOREIGN KEY ("id") REFERENCES "BookedVenues" ("venue_id");

ALTER TABLE "Profiles" ADD FOREIGN KEY ("id") REFERENCES "BookedVenues" ("booked_by");

ALTER TABLE "Profiles" ADD FOREIGN KEY ("id") REFERENCES "Favourites" ("added_by");

ALTER TABLE "Practitioners" ADD FOREIGN KEY ("id") REFERENCES "BookedPractitioners" ("service_id");

ALTER TABLE "Profiles" ADD FOREIGN KEY ("id") REFERENCES "BookedPractitioners" ("booked_by");

ALTER TABLE "Profiles" ADD FOREIGN KEY ("id") REFERENCES "Practitioners" ("created_by");

ALTER TABLE "Events" ADD FOREIGN KEY ("id") REFERENCES "Purchases" ("event_id");

ALTER TABLE "Venues" ADD FOREIGN KEY ("id") REFERENCES "Purchases" ("venue_id");

ALTER TABLE "Practitioners" ADD FOREIGN KEY ("id") REFERENCES "Purchases" ("service_id");

ALTER TABLE "Profiles" ADD FOREIGN KEY ("id") REFERENCES "Purchases" ("purchased_by");
