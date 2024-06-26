"use client";

import { useEffect, useReducer, useState } from "react";
import { useRouter } from "next/navigation";

type Task = {
    id: number;
    user_id: number;
    scheduled_time: string;
    task_type: string;
    task_details: {
        payload: {
            location: string;
            name: string;
            to_name: string;
            to: string;
        };
    };
};

function ScheduledTasks() {
    const [tasks, setTasks] = useState<Task[]>([]);
    const router = useRouter();

    useEffect(() => {
        const _body = {
            action: "fetch-tasks",
            empty: {},
        };
        async function fetch_tasks() {
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
                setTasks(result.data);
            } catch (error) {
                console.log(error);
            }
        }

        fetch_tasks();
    }, []);

    useEffect(() => {
        if (tasks) {
            console.log("tasks state updated: ", tasks);
        }
    }, [tasks]);

    return (
        <>
            <main className="flex flex-col w-full min-h-screen p-24">
                <div className="flex space-x-2 items-center">
          <span className="mb-2 text-gray-800 text-2xl font-bold mr-5">
            Scheduled Tasks
          </span>
                    <button
                        onClick={() => router.back()}
                        className="tracking-lighter hover:underline hover:text-blue-500 bg-gray-50 px-2 py-1"
                    >
                        back
                    </button>
                </div>
                <div className="flex flex-col space-y-5 w-full min-h-[800px] max-h-[800px] bg-gray-50 overflow-y-scroll p-12">
                    {tasks.map((task) => (
                        <>
                            <div className="flex w-full h-auto bg-white border border-gray-200 rounded-md p-5 shadow-lg">
                                <div className="flex flex-col space-y-1">
                  <span className="text-gray-500 text-md font-semibold tracking-wide">
                    ID:{task.id}
                  </span>
                                    <span className="text-gray-800 font-semibold text-xl">
                    Task Type: {task.task_type}
                  </span>
                                    <span className="text-sm text-gray-600 font-bold truncate">
                    receiver: {task.task_details.payload.to_name}
                  </span>
                                    <span className="text-sm text-gray-600 font-bold truncate">
                    Property Name: {task.task_details.payload.name}
                  </span>
                                    <span className="text-sm text-gray-600 font-bold truncate">
                    Property Location: {task.task_details.payload.location}
                  </span>
                                    <span className="text-sm text-gray-600 font-bold truncate">
                    Scheduled time: {task.scheduled_time}
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

export default ScheduledTasks;