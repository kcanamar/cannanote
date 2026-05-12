import 'package:flutter/material.dart';

import 'di/injection.dart';
import 'core/application/session_service.dart';
import 'core/domain/session.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize dependency injection
  await configureDependencies();

  runApp(const CannaNote());
}

class CannaNote extends StatelessWidget {
  const CannaNote({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'CannaNote',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: Colors.green,
          brightness: Brightness.light,
        ),
        useMaterial3: true,
      ),
      darkTheme: ThemeData(
        colorScheme: ColorScheme.fromSeed(
          seedColor: Colors.green,
          brightness: Brightness.dark,
        ),
        useMaterial3: true,
      ),
      home: const HomePage(),
    );
  }
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _sessionService = getIt<SessionService>();

  List<Session> _sessions = [];
  String _selectedMethod = 'vape';
  final _amountController = TextEditingController();
  String _statusMessage = '';

  final _methods = ['vape', 'smoke', 'edible', 'tincture'];

  @override
  void initState() {
    super.initState();
    _refreshSessions();
  }

  @override
  void dispose() {
    _amountController.dispose();
    super.dispose();
  }

  Future<void> _refreshSessions() async {
    final sessions = await _sessionService.getAllSessions();
    setState(() {
      _sessions = sessions;
    });
  }

  Future<void> _logSession() async {
    if (_amountController.text.isEmpty) {
      setState(() {
        _statusMessage = 'Please enter an amount';
      });
      return;
    }

    try {
      await _sessionService.logSession(
        method: _selectedMethod,
        amount: _amountController.text,
        unit: 'puffs',
      );

      setState(() {
        _statusMessage = 'Session logged!';
        _amountController.clear();
      });

      await _refreshSessions();
    } catch (e) {
      setState(() {
        _statusMessage = 'Error: $e';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('CannaNote'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Log Session Card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Log Session',
                      style: Theme.of(context).textTheme.titleMedium,
                    ),
                    const SizedBox(height: 16),

                    // Method dropdown
                    InputDecorator(
                      decoration: const InputDecoration(
                        labelText: 'Method',
                        border: OutlineInputBorder(),
                      ),
                      child: DropdownButtonHideUnderline(
                        child: DropdownButton<String>(
                          value: _selectedMethod,
                          isDense: true,
                          isExpanded: true,
                          items: _methods.map((method) {
                            return DropdownMenuItem(
                              value: method,
                              child: Text(method[0].toUpperCase() + method.substring(1)),
                            );
                          }).toList(),
                          onChanged: (value) {
                            setState(() {
                              _selectedMethod = value!;
                            });
                          },
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),

                    // Amount input
                    TextField(
                      controller: _amountController,
                      decoration: const InputDecoration(
                        labelText: 'Amount',
                        suffixText: 'puffs',
                        border: OutlineInputBorder(),
                      ),
                      keyboardType: TextInputType.number,
                    ),
                    const SizedBox(height: 16),

                    // Log button
                    FilledButton(
                      onPressed: _logSession,
                      child: const Text('Log Session'),
                    ),

                    // Status message
                    if (_statusMessage.isNotEmpty) ...[
                      const SizedBox(height: 8),
                      Text(
                        _statusMessage,
                        style: TextStyle(
                          color: _statusMessage.startsWith('Error')
                              ? Theme.of(context).colorScheme.error
                              : Theme.of(context).colorScheme.primary,
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Sessions List
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Sessions',
                  style: Theme.of(context).textTheme.titleMedium,
                ),
                IconButton(
                  icon: const Icon(Icons.refresh),
                  onPressed: _refreshSessions,
                ),
              ],
            ),
            const SizedBox(height: 8),

            Expanded(
              child: _sessions.isEmpty
                  ? const Center(
                      child: Text('No sessions yet. Log your first session above!'),
                    )
                  : ListView.builder(
                      itemCount: _sessions.length,
                      itemBuilder: (context, index) {
                        final session = _sessions[index];
                        return ListTile(
                          leading: const Icon(Icons.local_fire_department),
                          title: Text(
                            '${session.amount} ${session.unit} (${session.method})',
                          ),
                          subtitle: Text(
                            _formatTime(session.timestamp),
                          ),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  String _formatTime(DateTime dt) {
    return '${dt.hour.toString().padLeft(2, '0')}:${dt.minute.toString().padLeft(2, '0')}';
  }
}
