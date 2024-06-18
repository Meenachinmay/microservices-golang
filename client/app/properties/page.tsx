"use client";

import PropertyCard from "@/components/Property.card";
import { useSearchParams } from "next/navigation";

type Property = {
  uniqueId: string;
  name: string;
  location: string;
};

const properties: Property[] = [
  {
    uniqueId: "1",
    name: "Sakura Heights",
    location: "Tokyo, Shibuya",
  },
  {
    uniqueId: "2",
    name: "Hikari Tower",
    location: "Osaka, Namba",
  },
  {
    uniqueId: "3",
    name: "Fuji Gardens",
    location: "Kyoto, Fushimi",
  },
  {
    uniqueId: "4",
    name: "Kaze Villa",
    location: "Hokkaido, Sapporo",
  },
  {
    uniqueId: "5",
    name: "Yuki Residence",
    location: "Fukuoka, Hakata",
  },
  {
    uniqueId: "6",
    name: "Hana Apartments",
    location: "Nagoya, Naka",
  },
  {
    uniqueId: "7",
    name: "Umi House",
    location: "Kobe, Chuo",
  },
  {
    uniqueId: "8",
    name: "Tsubasa Estate",
    location: "Sendai, Aoba",
  },
  {
    uniqueId: "9",
    name: "Mizu Terrace",
    location: "Yokohama, Minato Mirai",
  },
  {
    uniqueId: "10",
    name: "Sakura Villa",
    location: "Hiroshima, Naka",
  },
];

type EnquiryPayload = {
  action: string;
  enquiry: {
    user_id: string | null;
    property_id: string;
    location: string
    name: string;
  };
};

function Properties() {
  const searchParams = useSearchParams();
  const userId = searchParams.get("userId");

  async function handleEnquiry(enquiry: Property) {
    const payload: EnquiryPayload = {
      action: "enquiry_mail",
      enquiry: {
        user_id: userId,
        property_id: enquiry.uniqueId,
        name: enquiry.name,
        location: enquiry.location,
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
