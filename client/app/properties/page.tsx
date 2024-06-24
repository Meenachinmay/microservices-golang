"use client";

import PropertyCard from "@/components/Property.card";
import { useSearchParams } from "next/navigation";
import {useEffect, useState} from "react";

type Property = {
  id: number
  name: string;
  location: string;
  fudousan_id: number
};


type EnquiryPayload = {
  user_id: number;
  property_id: number;
  property_location: string
  property_name: string
  available_timings: string
  preferred_method: string
  fudousan_id: number
};

type User = {
  userId: number,
  preferredMethod: string
  availableTimings: string
}

function Properties() {
  const [user, setUser] = useState<User>();
  const [properties, setProperties] = useState<Property[]>([]);

  async function handleEnquiry(enquiry: Property) {
    if (!user) {
     alert("user is not loaded");
     return
    }

    const payload: EnquiryPayload = {
        user_id: Number(user.userId),
        property_id: enquiry.id,
        property_name: enquiry.name,
        property_location: enquiry.location,
        available_timings:user.availableTimings,
        preferred_method:user.preferredMethod,
        fudousan_id: enquiry.fudousan_id,
    }

    try {
      const response = await fetch("http://localhost:8080/enquiry-grpc", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        console.error("error in response...");
        throw new Error("ERROR CALLING API");
      }

      const result = await response.json();
      console.log(result);
    } catch (error) {
      console.error(error);
      throw new Error("ERROR CALLING API");
    }
  }

  useEffect(() => {
    const tempUser = localStorage.getItem("current_user")
    setUser(JSON.parse(tempUser!));
    if (user) {
      console.log("user loaded: " + user);
    }

    async function fetchEnquiries(): Promise<void> {
      const _body = {
        action: "fetch-all-properties",
        empty: {}
      }
      try {
        const response = await fetch("http://localhost:8080/handle", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(_body)
        })

        if (!response.ok) {
          console.error("error in response...");
        }

        const result = await response.json();
        console.log(result);
        setProperties(result.data)
      } catch(error) {
        console.error(error);
        throw new Error("ERROR CALLING API");
      }
    }

    fetchEnquiries()

  }, [user?.userId]);

  return (
    <>
      <div className="flex items-center justify-center w-full h-[70px] bg-gray-800">
        <p className="text-center text-xl font-bold bg-gradient-to-r from-purple-500 via-pink-500 to-yellow-500 py-1 px-2 rounded-md">
          Enquiry for any property and we will contact you within 90 seconds.
        </p>
      </div>
      <main className="flex w-full min-h-screen bg-gradient-to-r from-gray-300 via-gray-200 to-gray-100 p-12">
        <div className="flex w-full min-h-full justify-center p-8">
          <div className="grid sm:grid-cols-1 md:grid-cols-4 gap-4">
            {properties.map((property, index) => (
              <PropertyCard
                key={index}
                name={property.name}
                location={property.location}
                handleEnquiry={() => handleEnquiry(property)}
                fudousan_id={property.fudousan_id}
              />
            ))}
          </div>
        </div>
      </main>
    </>
  );
}

export default Properties;
