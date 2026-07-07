import { pgTable, serial, integer, text, timestamp } from 'drizzle-orm/pg-core';

export const task = pgTable('task', {
	id: serial('id').primaryKey(),
	title: text('title').notNull(),
	priority: integer('priority').notNull().default(1)
});

/**
 * Application users for Auth.js credentials auth + Casbin authorization.
 * `role` is the Casbin subject (e.g. `admin`, `operator`, `viewer`).
 */
export const users = pgTable('users', {
	id: text('id')
		.primaryKey()
		.$defaultFn(() => crypto.randomUUID()),
	email: text('email').notNull().unique(),
	passwordHash: text('password_hash').notNull(),
	name: text('name').notNull(),
	role: text('role').notNull().default('viewer'),
	createdAt: timestamp('created_at').defaultNow().notNull()
});
