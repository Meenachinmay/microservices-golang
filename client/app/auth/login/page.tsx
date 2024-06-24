"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export type User = {
    userId: number,
    preferredMethod: string
    availableTimings: string
}

function Login() {
    const router = useRouter();
    const [userId, setUserId] = useState<string>("");
    const [preferredMethod, setPreferredMethod] = useState<string>("");
    const [availableTimings, setAvailableTimings] = useState<string>("");

    const handleSubmit = () => {
        let user: User = {
            userId: Number(userId),
            preferredMethod: preferredMethod,
            availableTimings: availableTimings
        }
        localStorage.setItem("current_user", JSON.stringify(user));
        router.push("/properties")
    }

    return (
        <>
            <main className="flex w-full min-h-screen bg-white">
                <div className="flex flex-col space-y-5 items-center justify-center w-1/2 h-[100vh] bg-gradient-to-t from-gray-50 via-gray-100 to-gray-200 p-12">
                    <input
                        onChange={(e) => setUserId(e.target.value)}
                        type="text"
                        placeholder="Enter UserId"
                        className="w-full h-[69px] px-4 border border-gray-300 focus:outline-none text-gray-800 font-semibold"
                        required
                    />
                    <input
                        onChange={(e) => setAvailableTimings(e.target.value)}
                        type="text"
                        placeholder="Enter available Time: HH:MM(13:00-15:00)"
                        className="w-full h-[69px] px-4 border border-gray-300 focus:outline-none text-gray-800 font-semibold"
                        required
                    />
                    <select
                        required
                        onChange={(e) => setPreferredMethod(e.target.value)}
                        className="w-full h-[69px] px-4 border border-gray-300 focus:outline-none text-gray-400 font-semibold">
                        <option value={""} disabled selected>
                            Enter preferred contact choice
                        </option>
                        <option value={"phone"}>Phone</option>
                        <option value={"sms"}>Sms</option>
                        <option value={"email"}>Email</option>
                    </select>
                    <button onClick={handleSubmit} className="w-full h-[69px] bg-orange-800 font-bold text-2xl hover:bg-orange-700">Go</button>
                </div>
                <div className="flex flex-col items-center justify-center w-1/2 h-[100vh] bg-orange-800 p-24">
                    <div className="flex flex-col">
            <span className="text-xl">
              Enter your userId to enter, so you can make enquires
            </span>
                        <span className="text-4xl font-semibold">
              We will reach you within 90 seconds anyhow.
            </span>
                    </div>
                </div>
            </main>
        </>
    );
}

export default Login;