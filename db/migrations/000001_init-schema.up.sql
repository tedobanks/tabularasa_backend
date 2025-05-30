CREATE TABLE "venues" (
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
  "owned_by" uuid, -- This is the foreign key column in 'venues'
  "is_available" boolean DEFAULT (false),
  "opens_at" timestamp,
  "closes_at" timestamp,
  "rental_days" varchar(255),
  "booking_price" integer,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "profiles" (
  "id" uuid PRIMARY KEY, -- This is referenced by other tables' FKs
  "bio" varchar(255),
  "phone_no" varchar(20),
  "country" varchar(255),
  "address" varchar(255),
  "experience" integer,
  "field" varchar(255),
  "business_name" varchar(255),
  "roles" VARCHAR(50) NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "practitioners" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "name" varchar(255) NOT NULL,
  "description" varchar(255) NOT NULL,
  "image_link" varchar(255),
  "is_available" boolean DEFAULT (true),
  "created_by" uuid, -- This is the foreign key column in 'practitioners'
  "opens_at" timestamp,
  "closes_at" timestamp,
  "working_days" varchar(255),
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "events" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "venue_id" uuid, -- This is the foreign key column in 'events'
  "image_links" varchar[],
  "name" varchar(255),
  "theme" varchar(255),
  "description" varchar(255),
  "audience" varchar(255),
  "activities" varchar[],
  "created_by" uuid, -- This is the foreign key column in 'events'
  "start_time" timestamp,
  "start_date" date,
  "end_date" date,
  "total_particpant" integer,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "favourites" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "event_id" uuid, -- This is the foreign key column in 'favourites'
  "added_by" uuid,   -- This is the foreign key column in 'favourites'
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()), -- This is referenced by profiles_users
  "email" varchar(255) UNIQUE NOT NULL,
  "password" varchar(255),
  "firstname" varchar(255),
  "lastname" varchar(255),
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "bookedVenues" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "type" varchar(255),
  "venue_id" uuid, -- This is the foreign key column in 'bookedVenues'
  "booked_for" timestamp,
  "booked_by" uuid, -- This is the foreign key column in 'bookedVenues'
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "bookedPractitioners" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "type" varchar(255),
  "service_id" uuid, -- This is the foreign key column in 'bookedPractitioners'
  "booked_for" timestamp,
  "booked_by" uuid, -- This is the foreign key column in 'bookedPractitioners'
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "purchases" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
  "event_id" uuid,    -- This is the foreign key column in 'purchases'
  "venue_id" uuid,    -- This is the foreign key column in 'purchases'
  "service_id" uuid,  -- This is the foreign key column in 'purchases'
  "purchased_by" uuid, -- This is the foreign key column in 'purchases'
  "created_at" timestamp DEFAULT (now())
);

-- Corrected Foreign Key Definitions:

-- A venue is owned by a profile
ALTER TABLE "venues" ADD FOREIGN KEY ("owned_by") REFERENCES "profiles" ("id");

-- An event is created by a profile
ALTER TABLE "events" ADD FOREIGN KEY ("created_by") REFERENCES "profiles" ("id");

-- An event is associated with a venue
ALTER TABLE "events" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");

-- Junction table for Profiles and Users (Many-to-Many)
CREATE TABLE "profiles_users" (
  "profiles_id" uuid,
  "users_id" uuid,
  PRIMARY KEY ("profiles_id", "users_id")
);

ALTER TABLE "profiles_users" ADD FOREIGN KEY ("profiles_id") REFERENCES "profiles" ("id");
ALTER TABLE "profiles_users" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id");

-- Favourites links events and profiles
ALTER TABLE "favourites" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");
ALTER TABLE "favourites" ADD FOREIGN KEY ("added_by") REFERENCES "profiles" ("id");

-- Booked Venues links venues and profiles
ALTER TABLE "bookedVenues" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");
ALTER TABLE "bookedVenues" ADD FOREIGN KEY ("booked_by") REFERENCES "profiles" ("id");

-- Booked Practitioners links practitioners (services) and profiles
ALTER TABLE "bookedPractitioners" ADD FOREIGN KEY ("service_id") REFERENCES "practitioners" ("id");
ALTER TABLE "bookedPractitioners" ADD FOREIGN KEY ("booked_by") REFERENCES "profiles" ("id");

-- Practitioners are created by profiles
ALTER TABLE "practitioners" ADD FOREIGN KEY ("created_by") REFERENCES "profiles" ("id");

-- Purchases link events, venues, services (practitioners), and profiles
ALTER TABLE "purchases" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");
ALTER TABLE "purchases" ADD FOREIGN KEY ("venue_id") REFERENCES "venues" ("id");
ALTER TABLE "purchases" ADD FOREIGN KEY ("service_id") REFERENCES "practitioners" ("id");
ALTER TABLE "purchases" ADD FOREIGN KEY ("purchased_by") REFERENCES "profiles" ("id");