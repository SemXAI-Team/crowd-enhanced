import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { users } from '$lib/server/db/schema';
import { eq } from 'drizzle-orm';
import bcrypt from 'bcryptjs';

export async function PUT({ params, request, locals }) {
	const session = await locals.auth();
	if (!session || session.user.role !== 'admin') {
		return json({ error: 'Unauthorized' }, { status: 401 });
	}

	try {
		const { password } = await request.json();
		if (!password) {
			return json({ error: 'Missing password' }, { status: 400 });
		}

		const salt = await bcrypt.genSalt(10);
		const passwordHash = await bcrypt.hash(password, salt);

		await db.update(users).set({ passwordHash }).where(eq(users.id, params.id));

		return json({ success: true });
	} catch (error) {
		console.error('Error updating password:', error);
		return json({ error: 'Failed to update password' }, { status: 500 });
	}
}
