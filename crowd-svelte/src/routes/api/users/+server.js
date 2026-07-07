import { json } from '@sveltejs/kit';
import { db } from '$lib/server/db';
import { users } from '$lib/server/db/schema';
import bcrypt from 'bcryptjs';

export async function POST({ request, locals }) {
	const session = await locals.auth();
	if (!session || session.user.role !== 'admin') {
		return json({ error: 'Unauthorized' }, { status: 401 });
	}

	try {
		const { name, email, password, role } = await request.json();
		if (!name || !email || !password || !role) {
			return json({ error: 'Missing required fields' }, { status: 400 });
		}

		const salt = await bcrypt.genSalt(10);
		const passwordHash = await bcrypt.hash(password, salt);

		await db.insert(users).values({
			name,
			email,
			passwordHash,
			role
		});

		return json({ success: true });
	} catch (error) {
		console.error('Error creating user:', error);
		return json({ error: 'Failed to create user' }, { status: 500 });
	}
}
