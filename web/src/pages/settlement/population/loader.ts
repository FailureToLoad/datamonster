import { redirect, type LoaderFunctionArgs } from "react-router";
import { Get } from "~/lib/request.ts";

export async function loadSurvivors({ params }: LoaderFunctionArgs) {
    const res = await Get(`/api/settlements/${params.settlementId}/survivors`);
    if (res.status === 401) return redirect('/');
    if (!res.ok) throw new Response('Failed to load survivors', { status: res.status });
    return res.json();
}