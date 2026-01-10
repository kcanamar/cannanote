-- Beta profiles table for production Supabase
CREATE TABLE profiles (
  id UUID REFERENCES auth.users(id) PRIMARY KEY,
  email TEXT NOT NULL,
  beta_status TEXT DEFAULT 'waitlist',
  beta_joined_at TIMESTAMPTZ DEFAULT NOW(),
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;

-- Allow anyone to create their own profile (during beta signup)
CREATE POLICY "Users can create their own profile" ON profiles
  FOR INSERT WITH CHECK (auth.uid() = id);

-- Users can view their own profile
CREATE POLICY "Users can view own profile" ON profiles
  FOR SELECT USING (auth.uid() = id);

-- Admin policy to view all profiles (for beta management)
CREATE POLICY "Service role can view all profiles" ON profiles
  FOR ALL USING (current_setting('role') = 'service_role');