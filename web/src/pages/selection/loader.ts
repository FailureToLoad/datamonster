import { redirect } from "react-router";
import { Get } from "~/lib/request.ts";

export async function loadSettlements() {
    const res = await Get("/api/settlements");
    if (res.status === 401) return redirect("/");
    if (!res.ok) throw new Response("Failed to load", { status: res.status });
    return res.json();
}
