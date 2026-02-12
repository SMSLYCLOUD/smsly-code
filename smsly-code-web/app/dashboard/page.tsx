"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

interface Repo {
  id: number;
  name: string;
  description: string;
  is_private: boolean;
}

export default function Dashboard() {
  const router = useRouter();
  const [repos, setRepos] = useState<Repo[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/login");
      return;
    }

    fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/repos`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch repos");
        return res.json();
      })
      .then((data) => {
        setRepos(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error(err);
        setLoading(false);
      });
  }, [router]);

  return (
    <div className="min-h-screen p-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Your Repositories</h1>
          <Link
            href="/repos/new"
            className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-500"
          >
            New Repository
          </Link>
        </div>

        {loading ? (
          <p>Loading...</p>
        ) : repos.length === 0 ? (
          <div className="text-center py-12 bg-gray-50 rounded-lg border border-gray-200">
            <p className="text-gray-500 mb-4">You don't have any repositories yet.</p>
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {repos.map((repo) => (
              <Link
                href={`/repos/${repo.name}`}
                key={repo.id}
                className="block p-6 bg-white border border-gray-200 rounded-lg hover:border-indigo-500 transition-colors"
              >
                <div className="flex justify-between items-start mb-2">
                  <h2 className="text-xl font-semibold text-gray-900">
                    {repo.name}
                  </h2>
                  <span
                    className={`px-2 py-1 text-xs rounded-full ${
                      repo.is_private
                        ? "bg-gray-100 text-gray-800"
                        : "bg-green-100 text-green-800"
                    }`}
                  >
                    {repo.is_private ? "Private" : "Public"}
                  </span>
                </div>
                <p className="text-gray-600 line-clamp-2">
                  {repo.description || "No description provided."}
                </p>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
