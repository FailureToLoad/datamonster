import type { Survivor } from "~/lib/survivor";
import type { SettlementId } from "~/lib/settlement";

export type MockApiState = {
  authenticated: boolean;
  settlements: SettlementId[];
  survivors: Map<string, Survivor[]>;
  glossary: {
    disorders: Array<{ id: string; name: string }>;
    fightingArts: Array<{ id: string; name: string }>;
  };
};

export function createMockApiState(
  overrides: Partial<MockApiState> = {}
): MockApiState {
  return {
    authenticated: true,
    settlements: [],
    survivors: new Map(),
    glossary: { disorders: [], fightingArts: [] },
    ...overrides,
  };
}

export function createMockFetch(state: MockApiState) {
  return async (
    input: RequestInfo | URL,
    init?: RequestInit
  ): Promise<Response> => {
    const url = typeof input === "string" ? input : input.toString();
    const method = init?.method?.toUpperCase() ?? "GET";

    if (!state.authenticated && !url.includes("/api/auth")) {
      return new Response(null, { status: 401 });
    }

    if (url === "/api/me") {
      return state.authenticated
        ? Response.json({ id: "user-1", email: "test@example.com" })
        : new Response(null, { status: 401 });
    }

    if (url === "/api/auth/logout" && method === "POST") {
      state.authenticated = false;
      return new Response(null, { status: 200 });
    }

    if (url === "/api/glossary") {
      return Response.json(state.glossary);
    }

    if (url === "/api/settlements" && method === "GET") {
      return Response.json(state.settlements);
    }

    if (url === "/api/settlements" && method === "POST") {
      const body = JSON.parse(init?.body as string);
      const newSettlement: SettlementId = {
        id: `settlement-${Date.now()}`,
        name: body.name,
      };
      state.settlements.push(newSettlement);
      return Response.json(newSettlement, { status: 201 });
    }

    const survivorMatch = url.match(/\/api\/settlements\/([^/]+)\/survivors$/);
    if (survivorMatch) {
      const settlementId = survivorMatch[1];

      if (method === "GET") {
        return Response.json(state.survivors.get(settlementId) ?? []);
      }

      if (method === "POST") {
        const body = JSON.parse(init?.body as string);
        const newSurvivor: Survivor = {
          id: `survivor-${Date.now()}`,
          settlementId,
          name: body.name,
          gender: body.gender ?? "M",
          birth: body.birth ?? 1,
          status: body.status ?? "Alive",
          huntxp: body.huntxp ?? 0,
          movement: body.movement ?? 5,
          speed: body.speed ?? 0,
          strength: body.strength ?? 0,
          accuracy: body.accuracy ?? 0,
          evasion: body.evasion ?? 0,
          luck: body.luck ?? 0,
          systemicPressure: body.systemicPressure ?? 0,
          torment: body.torment ?? 0,
          courage: body.courage ?? 0,
          understanding: body.understanding ?? 0,
          survival: body.survival ?? 0,
          insanity: body.insanity ?? 0,
          lumi: body.lumi ?? 0,
          disorders: body.disorders ?? [],
          fightingArt: body.fightingArt ?? null,
          secretFightingArt: body.secretFightingArt ?? null,
        };
        const existing = state.survivors.get(settlementId) ?? [];
        state.survivors.set(settlementId, [...existing, newSurvivor]);
        return Response.json(newSurvivor, { status: 201 });
      }
    }

    const survivorPatchMatch = url.match(
      /\/api\/settlements\/([^/]+)\/survivors\/([^/]+)$/
    );
    if (survivorPatchMatch && method === "PATCH") {
      const [, settlementId, survivorId] = survivorPatchMatch;
      const survivors = state.survivors.get(settlementId) ?? [];
      const index = survivors.findIndex((s) => s.id === survivorId);
      if (index === -1) {
        return new Response(null, { status: 404 });
      }
      const body = JSON.parse(init?.body as string);
      const updated = { ...survivors[index], ...body.statUpdates };
      if (body.statusUpdate) updated.status = body.statusUpdate;
      if (body.disorders) updated.disorders = body.disorders;
      if (body.fightingArt !== undefined)
        updated.fightingArt = body.fightingArt;
      if (body.secretFightingArt !== undefined)
        updated.secretFightingArt = body.secretFightingArt;
      survivors[index] = updated;
      return Response.json(updated);
    }

    return new Response(null, { status: 404 });
  };
}
