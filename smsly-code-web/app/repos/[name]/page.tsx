"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Issue } from "../../../types";

interface Commit {
  id: string;
  message: string;
  author: string;
  date: string;
  mip_verified: boolean;
}

interface TreeEntry {
  name: string;
  kind: string;
  id: string;
}

export default function RepoView({ params }: { params: { name: string } }) {
  const router = useRouter();
  const [commits, setCommits] = useState<Commit[]>([]);
  const [files, setFiles] = useState<TreeEntry[]>([]);
  const [issues, setIssues] = useState<Issue[]>([]);
  const [activeTab, setActiveTab] = useState<"code" | "commits" | "issues">("code");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/login");
      return;
    }

    const fetchData = async () => {
      try {
        const headers = { Authorization: `Bearer ${token}` };

        const filesRes = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/repos/${params.name}/tree/HEAD`,
          { headers }
        );
        if (filesRes.ok) {
          const filesData = await filesRes.json();
          if (Array.isArray(filesData)) {
            setFiles(filesData);
          }
        }

        const commitsRes = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/repos/${params.name}/commits`,
          { headers }
        );
        if (commitsRes.ok) {
          const commitsData = await commitsRes.json();
          if (Array.isArray(commitsData)) {
             setCommits(commitsData);
          }
        }

        const issuesRes = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/repos/${params.name}/issues`,
          { headers }
        );
        if (issuesRes.ok) {
          const issuesData = await issuesRes.json();
          if (Array.isArray(issuesData)) {
             setIssues(issuesData);
          }
        }

      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [params.name, router]);

  if (loading) return <div className="p-8">Loading repository data...</div>;

  return (
    <div className="min-h-screen p-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-3xl font-bold">{params.name}</h1>
          <div className="flex gap-2">
             <span className="px-3 py-1 bg-gray-100 text-gray-800 rounded-full text-sm font-medium border border-gray-200">
               Public
             </span>
             <Link href={`/repos/${params.name}/issues/new`} className="px-3 py-1 bg-indigo-600 text-white rounded-full text-sm font-medium hover:bg-indigo-700">
               New Issue
             </Link>
          </div>
        </div>

        <div className="mb-6">
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex space-x-8" aria-label="Tabs">
              <button
                onClick={() => setActiveTab("code")}
                className={`${
                  activeTab === "code"
                    ? "border-indigo-500 text-indigo-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm`}
              >
                Code
              </button>
              <button
                onClick={() => setActiveTab("commits")}
                className={`${
                  activeTab === "commits"
                    ? "border-indigo-500 text-indigo-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm`}
              >
                Commits
              </button>
              <button
                onClick={() => setActiveTab("issues")}
                className={`${
                  activeTab === "issues"
                    ? "border-indigo-500 text-indigo-600"
                    : "border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300"
                } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm`}
              >
                Issues
              </button>
            </nav>
          </div>
        </div>

        {activeTab === "code" && (
          <div className="bg-white shadow overflow-hidden sm:rounded-md border border-gray-200">
            <div className="bg-gray-50 px-4 py-3 border-b border-gray-200 sm:px-6">
               <h3 className="text-lg leading-6 font-medium text-gray-900">Files</h3>
            </div>
            {files.length === 0 ? (
              <div className="p-8 text-center text-gray-500">
                <p className="mb-4">This repository is empty.</p>
                <div className="p-4 bg-gray-100 rounded text-left font-mono text-sm inline-block">
                   <p className="mb-1"># Push an existing repository</p>
                   <p>git remote add origin http://localhost:8081/{params.name}</p>
                   <p>git push -u origin main</p>
                </div>
              </div>
            ) : (
              <ul className="divide-y divide-gray-200">
                {files.map((file) => (
                  <li key={file.id}>
                    <div className="px-4 py-4 flex items-center sm:px-6 hover:bg-gray-50 cursor-pointer transition">
                      <div className="min-w-0 flex-1 sm:flex sm:items-center sm:justify-between">
                        <div className="flex items-center truncate">
                          <span className="mr-3 text-xl">
                            {file.kind === "tree" ? "📁" : "📄"}
                          </span>
                          <p className="font-medium text-indigo-600 truncate">{file.name}</p>
                        </div>
                        <div className="mt-4 flex-shrink-0 sm:mt-0 sm:ml-5">
                          <div className="flex overflow-hidden -space-x-1">
                             <span className="text-xs text-gray-400 font-mono">{file.id.substring(0, 7)}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        )}

        {activeTab === "commits" && (
          <div className="bg-white shadow overflow-hidden sm:rounded-md border border-gray-200">
             <div className="bg-gray-50 px-4 py-3 border-b border-gray-200 sm:px-6">
               <h3 className="text-lg leading-6 font-medium text-gray-900">Commit History</h3>
            </div>
            <ul className="divide-y divide-gray-200">
              {commits.length === 0 ? (
                  <div className="p-8 text-center text-gray-500">No commits found.</div>
              ) : commits.map((commit) => (
                <li key={commit.id}>
                  <div className="px-4 py-4 sm:px-6 hover:bg-gray-50 transition">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center">
                        <p className="text-sm font-medium text-indigo-600 truncate mr-3">
                          {commit.message}
                        </p>
                        {commit.mip_verified && (
                          <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                            <svg className="mr-1.5 h-2 w-2 text-green-400" fill="currentColor" viewBox="0 0 8 8">
                              <circle cx="4" cy="4" r="3" />
                            </svg>
                            MIP Verified
                          </span>
                        )}
                      </div>
                      <div className="ml-2 flex-shrink-0 flex">
                        <p className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800 font-mono border border-gray-200">
                          {commit.id.substring(0, 7)}
                        </p>
                      </div>
                    </div>
                    <div className="mt-2 sm:flex sm:justify-between">
                      <div className="sm:flex">
                        <p className="flex items-center text-sm text-gray-500">
                          {commit.author}
                        </p>
                      </div>
                      <div className="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                        <p>
                          Committed on <time dateTime={commit.date}>{new Date(parseInt(commit.date) * 1000).toLocaleDateString()}</time>
                        </p>
                      </div>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        )}

        {activeTab === "issues" && (
          <div className="bg-white shadow overflow-hidden sm:rounded-md border border-gray-200">
             <div className="bg-gray-50 px-4 py-3 border-b border-gray-200 sm:px-6 flex justify-between items-center">
               <h3 className="text-lg leading-6 font-medium text-gray-900">Issues</h3>
            </div>
            <ul className="divide-y divide-gray-200">
              {issues.length === 0 ? (
                  <div className="p-8 text-center text-gray-500">No issues found.</div>
              ) : issues.map((issue) => (
                <li key={issue.id}>
                   <Link href={`/repos/${params.name}/issues/${issue.id}`} className="block hover:bg-gray-50 transition">
                    <div className="px-4 py-4 sm:px-6">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center">
                          {issue.state === "open" ? (
                            <span className="text-green-600 mr-2">⊙</span>
                          ) : (
                            <span className="text-red-600 mr-2">×</span>
                          )}
                          <p className="text-sm font-medium text-indigo-600 truncate mr-3">
                            {issue.title}
                          </p>
                        </div>
                        <div className="ml-2 flex-shrink-0 flex">
                          <p className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-gray-100 text-gray-800 font-mono border border-gray-200">
                            #{issue.id}
                          </p>
                        </div>
                      </div>
                      <div className="mt-2 sm:flex sm:justify-between">
                        <div className="sm:flex">
                          <p className="flex items-center text-sm text-gray-500">
                            Opened by {issue.creator?.username || "Unknown"}
                          </p>
                        </div>
                        <div className="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                          <p>
                            {new Date(issue.created_at).toLocaleDateString()}
                          </p>
                        </div>
                      </div>
                    </div>
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
