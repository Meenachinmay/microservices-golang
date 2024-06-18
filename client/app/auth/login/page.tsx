"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

function Login() {
    const router = useRouter();
    const [userId, setUserId] = useState<string>("");

    function handleChange (e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === 'Enter') {
        router.push(`/properties?userId=${userId}`)
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
        <div className="flex items-center justify-center w-full h-[200px] border border-gray-300">
          <input
            type="text"
            placeholder="userId"
            className="w-1/5 p-2 focus:outline-none border border-gray-300 text-gray-800 font-semibold"
            onChange={(e) => setUserId(e.target.value)}
            onKeyDown={handleChange} 
          />
        </div>
      </main>
    </>
  );
}

export default Login;
