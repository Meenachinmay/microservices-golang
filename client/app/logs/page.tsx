"use client";

import {useEffect, useState} from "react";
import {useRouter} from "next/navigation";

type Log = {
    id: number;
    service_name: string;
    log_data: string;
};

function Page() {
    const [logs, setLogs] = useState<Log[]>([]);
    const router = useRouter();

    useEffect(() => {
        const _body = {
            action: "fetch-log",
            empty: {},
        };

        async function fetch_logs() {
            try {
                const response = await fetch("http://localhost:8080/handle", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(_body),
                });

                if (!response.ok) {
                    console.error("ERROR CALLING API...");
                }

                const result = await response.json();
                setLogs(result.data);
            } catch (error) {
                console.log(error);
            }
        }

        fetch_logs();
    }, []);

    useEffect(() => {
        if (logs) {
            console.log("logs state updated: ", logs);
        }
    }, [logs]);

    return (
        <>
            <main className="flex flex-col w-full min-h-screen p-24">
                <div className="flex space-x-2 items-center">
          <span className="mb-2 text-gray-800 text-2xl font-bold mr-5">
            Logs
          </span>
                    <button
                        onClick={() => router.back()}
                        className="tracking-lighter hover:underline hover:text-blue-500 bg-gray-50 px-2 py-1"
                    >
                        back
                    </button>
                </div>
                <div
                    className="flex flex-col space-y-5 w-full max-h-[800px] min-h-full bg-gray-50 overflow-y-scroll p-12">
                    {logs.map((log) => (
                        <>
                            <div
                                className="flex w-full h-[125px] bg-white border border-gray-200 rounded-md p-5 shadow-lg">
                                <div className="flex flex-col space-y-1">
              <span className="text-gray-500 text-md font-semibold tracking-wide">
                ID: {log.id}
              </span>
                                    <span className="text-gray-800 font-semibold text-xl">
                Service Name: {log.service_name}
              </span>
                                    <span className="text-sm text-gray-600 font-bold truncate">
                Log Data: {log.log_data}
              </span>
                                </div>
                            </div>
                        </>
                    ))}
                </div>
            </main>
        </>
    );
}

export default Page;