"use client";

import { useState, useEffect, useCallback } from "react";
import ReactMarkdown from "react-markdown";
import { useRouter } from "next/navigation";
import { Issue, Comment } from "../../../../../types";

export default function IssueDetail({ params }: { params: { name: string; id: string } }) {
  const router = useRouter();
  const [issue, setIssue] = useState<Issue | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState("");
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async () => {
    const token = localStorage.getItem("token");
    if (!token) {
      router.push("/login");
      return;
    }

    try {
      const headers = { Authorization: `Bearer ${token}` };

      const issueRes = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/repos/${params.name}/issues/${params.id}`,
        { headers }
      );
      if (issueRes.ok) {
        setIssue(await issueRes.json());
      }

      const commentsRes = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/repos/${params.name}/issues/${params.id}/comments`,
        { headers }
      );
      if (commentsRes.ok) {
        setComments(await commentsRes.json());
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [params.name, params.id, router]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleAddComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim()) return;

    const token = localStorage.getItem("token");
    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/repos/${params.name}/issues/${params.id}/comments`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ body: newComment }),
        }
      );

      if (res.ok) {
        setNewComment("");
        fetchData(); // Reload comments
      }
    } catch (err) {
      console.error(err);
    }
  };

  const handleToggleState = async () => {
    if (!issue) return;
    const newState = issue.state === "open" ? "closed" : "open";

    const token = localStorage.getItem("token");
    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/repos/${params.name}/issues/${params.id}`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ state: newState }),
        }
      );

      if (res.ok) {
        setIssue({ ...issue, state: newState });
      }
    } catch (err) {
      console.error(err);
    }
  };

  if (loading) return <div className="p-8">Loading issue...</div>;
  if (!issue) return <div className="p-8">Issue not found</div>;

  return (
    <div className="min-h-screen p-8 bg-gray-50">
      <div className="max-w-4xl mx-auto">
        <div className="mb-6 flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold mb-2">
              {issue.title} <span className="text-gray-400">#{issue.id}</span>
            </h1>
            <div className="flex items-center gap-2">
              <span
                className={`px-3 py-1 rounded-full text-white text-sm font-medium ${
                  issue.state === "open" ? "bg-green-600" : "bg-red-600"
                }`}
              >
                {issue.state === "open" ? "Open" : "Closed"}
              </span>
              <span className="text-gray-600">
                {issue.creator?.username} opened this issue on {new Date(issue.created_at).toLocaleDateString()}
              </span>
            </div>
          </div>
          <button
            onClick={handleToggleState}
            className="px-4 py-2 border border-gray-300 rounded text-sm font-medium hover:bg-gray-100"
          >
            {issue.state === "open" ? "Close Issue" : "Reopen Issue"}
          </button>
        </div>

        <div className="bg-white border border-gray-200 rounded mb-6">
          <div className="p-4 border-b border-gray-100 bg-gray-50 flex justify-between">
            <span className="font-bold">{issue.creator?.username}</span>
            <span className="text-gray-500 text-sm">commented</span>
          </div>
          <div className="p-4 prose max-w-none">
            <ReactMarkdown>{issue.body}</ReactMarkdown>
          </div>
        </div>

        <div className="space-y-6 mb-8">
          {comments.map((comment) => (
            <div key={comment.id} className="bg-white border border-gray-200 rounded">
              <div className="p-4 border-b border-gray-100 bg-gray-50 flex justify-between">
                <span className="font-bold">{comment.user?.username}</span>
                <span className="text-gray-500 text-sm">
                   commented on {new Date(comment.created_at).toLocaleDateString()}
                </span>
              </div>
              <div className="p-4 prose max-w-none">
                <ReactMarkdown>{comment.body}</ReactMarkdown>
              </div>
            </div>
          ))}
        </div>

        <div className="bg-white border border-gray-200 rounded p-4">
          <h3 className="text-lg font-medium mb-4">Add a comment</h3>
          <form onSubmit={handleAddComment}>
            <textarea
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
              className="w-full p-3 border rounded mb-4 focus:ring-2 focus:ring-indigo-500 focus:outline-none"
              rows={4}
              placeholder="Leave a comment"
            />
            <div className="flex justify-end">
              <button
                type="submit"
                disabled={!newComment.trim()}
                className="px-4 py-2 bg-green-600 text-white rounded font-medium hover:bg-green-700 disabled:opacity-50"
              >
                Comment
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
