"use client";

type PropertyUICardProps = {
  name: string;
  location: string;
  handleEnquiry: () => void
};

type Property = {
  uniqueId: string;
  name: string;
  location: string;
};

function PropertyCard({ name, location, handleEnquiry }: PropertyUICardProps) {
  return (
    <>
      <div className="flex flex-col w-[350px] h-[250px] bg-white border border-gray-300
         hover:shadow-xl rounded-md cursor-pointer">
        <div className="flex flex-col w-full min-h-[70%] p-4">
          <p className="text-gray-800 font-bold text-2xl">{name}</p>
          <p className="text-gray-500 font-semibold text-sm">{location}</p>
        </div>
        <div className="flex justify-end w-full min-h-[30%] bg-gray-800 p-4">
          <button onClick={handleEnquiry} className="focus:outline-none text-sm bg-gradient-to-r from-blue-500 via-blue-600 to-blue-800 w-auto font-bold px-2 text-white">
            お問い合わせ
          </button>
        </div>
      </div>
    </>
  );
}

export default PropertyCard;
