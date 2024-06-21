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
  action: string;
  enquiry: {
    user_id: number;
    property_id: number;
    location: string
    name: string;
    fudousan_id: number
  };
};

function Properties() {
  const searchParams = useSearchParams();
  const userId = searchParams.get("userId");

  const [properties, setProperties] = useState<Property[]>([]);

  async function handleEnquiry(enquiry: Property) {
    const payload: EnquiryPayload = {
      action: "add_new_enquiry",
      enquiry: {
        user_id: enquiry.id,
        property_id: enquiry.id,
        name: enquiry.name,
        location: enquiry.location,
        fudousan_id: enquiry.fudousan_id,
      },
    };

    try {
      const response = await fetch("http://localhost:8080/handle", {
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
        setProperties(result.data)
      } catch(error) {
        console.error(error);
        throw new Error("ERROR CALLING API");
      }
    }

    fetchEnquiries()

  }, [userId]);

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
              />
            ))}
          </div>
        </div>
      </main>
    </>
  );
}

export default Properties;
